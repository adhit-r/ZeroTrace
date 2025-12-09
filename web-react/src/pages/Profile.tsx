import React from 'react';
import { User, Mail, Shield, Clock, Activity, Award } from 'lucide-react';

const Profile: React.FC = () => {
  return (
    <div className="p-6 bg-gray-100 min-h-screen">
      <div className="max-w-4xl mx-auto">
        <div className="mb-8">
          <h1 className="text-4xl font-black text-black uppercase tracking-wider mb-2">USER PROFILE</h1>
          <p className="text-lg text-gray-600 font-bold">MANAGE YOUR ACCOUNT AND PREFERENCES</p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Profile Card */}
          <div className="lg:col-span-1">
            <div className="bg-white border-3 border-black rounded-lg shadow-neubrutalist p-6">
              <div className="text-center mb-6">
                <div className="w-24 h-24 bg-orange-500 border-3 border-black rounded-full mx-auto mb-4 flex items-center justify-center">
                  <User className="w-12 h-12 text-white" />
                </div>
                <h2 className="text-2xl font-black text-black uppercase tracking-wider">JOHN DOE</h2>
                <p className="text-orange-500 font-bold">SECURITY ANALYST</p>
              </div>
              
              <div className="space-y-4">
                <div className="flex items-center">
                  <Mail className="w-5 h-5 text-orange-500 mr-3" />
                  <span className="font-bold text-black">john.doe@company.com</span>
                </div>
                <div className="flex items-center">
                  <Shield className="w-5 h-5 text-orange-500 mr-3" />
                  <span className="font-bold text-black">Level 3 Access</span>
                </div>
                <div className="flex items-center">
                  <Clock className="w-5 h-5 text-orange-500 mr-3" />
                  <span className="font-bold text-black">Last login: 2 hours ago</span>
                </div>
              </div>
            </div>
          </div>

          {/* Activity & Stats */}
          <div className="lg:col-span-2">
            <div className="space-y-6">
              {/* Activity Stats */}
              <div className="bg-white border-3 border-black rounded-lg shadow-neubrutalist p-6">
                <div className="flex items-center mb-4">
                  <Activity className="w-8 h-8 text-orange-500 mr-3" />
                  <h3 className="text-xl font-black text-black uppercase tracking-wider">ACTIVITY STATS</h3>
                </div>
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                  <div className="text-center">
                    <div className="text-3xl font-black text-orange-500 mb-2">247</div>
                    <div className="text-sm font-bold text-black uppercase">SCANS RUN</div>
                  </div>
                  <div className="text-center">
                    <div className="text-3xl font-black text-green-500 mb-2">1,892</div>
                    <div className="text-sm font-bold text-black uppercase">VULNERABILITIES FOUND</div>
                  </div>
                  <div className="text-center">
                    <div className="text-3xl font-black text-blue-500 mb-2">156</div>
                    <div className="text-sm font-bold text-black uppercase">REPORTS GENERATED</div>
                  </div>
                  <div className="text-center">
                    <div className="text-3xl font-black text-purple-500 mb-2">89%</div>
                    <div className="text-sm font-bold text-black uppercase">SUCCESS RATE</div>
                  </div>
                </div>
              </div>

              {/* Recent Activity */}
              <div className="bg-white border-3 border-black rounded-lg shadow-neubrutalist p-6">
                <div className="flex items-center mb-4">
                  <Clock className="w-8 h-8 text-orange-500 mr-3" />
                  <h3 className="text-xl font-black text-black uppercase tracking-wider">RECENT ACTIVITY</h3>
                </div>
                <div className="space-y-3">
                  <div className="flex items-center justify-between p-3 bg-gray-50 border-2 border-black rounded">
                    <div className="flex items-center">
                      <div className="w-3 h-3 bg-green-500 border-2 border-black rounded-full mr-3"></div>
                      <span className="font-bold text-black">Completed vulnerability scan</span>
                    </div>
                    <span className="text-sm text-gray-600">2 hours ago</span>
                  </div>
                  <div className="flex items-center justify-between p-3 bg-gray-50 border-2 border-black rounded">
                    <div className="flex items-center">
                      <div className="w-3 h-3 bg-orange-500 border-2 border-black rounded-full mr-3"></div>
                      <span className="font-bold text-black">Generated security report</span>
                    </div>
                    <span className="text-sm text-gray-600">5 hours ago</span>
                  </div>
                  <div className="flex items-center justify-between p-3 bg-gray-50 border-2 border-black rounded">
                    <div className="flex items-center">
                      <div className="w-3 h-3 bg-blue-500 border-2 border-black rounded-full mr-3"></div>
                      <span className="font-bold text-black">Updated agent configuration</span>
                    </div>
                    <span className="text-sm text-gray-600">1 day ago</span>
                  </div>
                </div>
              </div>

              {/* Achievements */}
              <div className="bg-white border-3 border-black rounded-lg shadow-neubrutalist p-6">
                <div className="flex items-center mb-4">
                  <Award className="w-8 h-8 text-orange-500 mr-3" />
                  <h3 className="text-xl font-black text-black uppercase tracking-wider">ACHIEVEMENTS</h3>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div className="text-center p-4 bg-orange-100 border-2 border-orange-300 rounded">
                    <div className="text-2xl mb-2"></div>
                    <div className="font-bold text-black uppercase">SCAN MASTER</div>
                    <div className="text-sm text-gray-600">100+ scans completed</div>
                  </div>
                  <div className="text-center p-4 bg-green-100 border-2 border-green-300 rounded">
                    <div className="text-2xl mb-2">️</div>
                    <div className="font-bold text-black uppercase">SECURITY EXPERT</div>
                    <div className="text-sm text-gray-600">500+ vulnerabilities found</div>
                  </div>
                  <div className="text-center p-4 bg-blue-100 border-2 border-blue-300 rounded">
                    <div className="text-2xl mb-2"></div>
                    <div className="font-bold text-black uppercase">REPORT GENIUS</div>
                    <div className="text-sm text-gray-600">50+ reports generated</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Profile;